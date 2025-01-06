import requests
import os
import argparse
import shutil
from selectolax.parser import HTMLParser
from download_utils import download_images_parallel
from playwright.sync_api import sync_playwright


def search_manga(title):
    url = "https://weebcentral.com/search/simple"
    querystring = {"location": "main"}
    payload = f"text={title}"
    headers = {
        "content-type": "application/x-www-form-urlencoded",
        "user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0",
    }

    response = requests.request(
        "POST", url, data=payload, headers=headers, params=querystring
    )

    # Parse HTML response
    parser = HTMLParser(response.text)

    # Find all manga links in search results
    search_results = parser.css("#quick-search-result a")

    if search_results:
        # Get the first result
        first_result = search_results[0]
        manga_url = first_result.attributes.get("href")
        manga_title = first_result.css_first("div.flex-1").text().strip()
        print(f"Found manga: {manga_title}")
        print(f"URL: {manga_url}")
        return manga_url

    print("No manga found in search results")
    return None


def extract_series_id(manga_url):
    # Extract the series ID from the URL
    # Example: from https://weebcentral.com/series/01J76XYFM1TWGNNQ2Y2T8V7E8Y/Wistoria-Wand-and-Sword
    # get 01J76XYFM1TWGNNQ2Y2T8V7E8Y
    parts = manga_url.split("/")
    for part in parts:
        if len(part) == 26 and part.isalnum():  # Series IDs are 26 characters long
            return part
    return None


def get_manga_slug(manga_url):
    # Extract the manga title slug from the URL
    parts = manga_url.split("/")
    if len(parts) >= 6:  # URL format: https://weebcentral.com/series/ID/SLUG
        return parts[-1]
    return None


def get_base_url(manga_url):
    # Get the base URL without the title slug
    series_id = extract_series_id(manga_url)
    if series_id:
        return f"https://weebcentral.com/series/{series_id}/"
    return None


def get_rss_url(manga_url):
    # Construct the RSS feed URL using the series ID
    series_id = extract_series_id(manga_url)
    if series_id:
        return f"https://weebcentral.com/series/{series_id}/rss"
    return None


def get_chapter_list_url(manga_url):
    # Construct the full chapter list URL using the base URL
    base_url = get_base_url(manga_url)
    if base_url:
        return f"{base_url}full-chapter-list"
    return None


def get_chapters_from_list(chapter_list_url):
    """Get chapter links from the full chapter list page"""
    headers = {
        "user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0",
    }

    chapters = []
    try:
        vprint(f"\nDEBUG: Fetching chapter list from URL: {chapter_list_url}")
        response = requests.get(chapter_list_url, headers=headers)
        response.raise_for_status()

        parser = HTMLParser(response.text)
        chapter_links = parser.css("a[href*='/chapters/']")
        vprint(f"DEBUG: Found {len(chapter_links)} raw chapter links on page")

        for link in chapter_links:
            # Get raw title and clean it
            title = link.text()
            # Clean up the title by removing extra whitespace and CSS/HTML content
            cleaned_title = " ".join(
                [
                    line.strip()
                    for line in title.split("\n")
                    if line.strip()
                    and not line.strip().startswith(".")  # Skip CSS
                    and not line.strip().startswith("{")  # Skip CSS
                    and not line.strip().startswith("2024")  # Skip dates
                ]
            ).strip()

            url = link.attributes.get("href")
            vprint(f"\nDEBUG: Processing link - Cleaned Title: {cleaned_title}")
            vprint(f"DEBUG: Link URL: {url}")

            # Extract chapter number from cleaned title
            chapter_num = None
            if "Chapter" in cleaned_title:
                parts = cleaned_title.split("Chapter")
                format_type = "Chapter"
            elif "Episode" in cleaned_title:
                parts = cleaned_title.split("Episode")
                format_type = "Episode"
            elif "Days" in cleaned_title:
                parts = cleaned_title.split("Days")
                format_type = "Days"
            else:
                parts = None
                format_type = None
                vprint("DEBUG: No recognized chapter format in title")

            if parts and len(parts) > 1:
                num_part = parts[-1].strip()
                vprint(
                    f"DEBUG: Found '{format_type}' format - Extracted number part: '{num_part}'"
                )
                # Only use the first number found
                num_part = num_part.split()[0]
                if num_part.replace(".", "").isdigit():
                    chapter_num = num_part
                    vprint(f"DEBUG: Valid chapter number found: {chapter_num}")

            if chapter_num:
                chapters.append({"chapter": chapter_num, "url": url})
                vprint(f"DEBUG: Added chapter {chapter_num} from {url}")

        # Sort chapters numerically
        sorted_chapters = sorted(
            chapters,
            key=lambda x: float(x["chapter"])
            if x["chapter"].replace(".", "").isdigit()
            else 0,
        )
        vprint(f"\nDEBUG: Final chapter count from list: {len(sorted_chapters)}")
        return sorted_chapters

    except Exception as e:
        print(f"Error getting chapter list: {str(e)}")
        return []


def get_chapter_links(rss_url):
    import xmltodict

    headers = {
        "user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0",
    }

    response = requests.get(rss_url, headers=headers)
    response.raise_for_status()

    # Parse XML response
    rss_data = xmltodict.parse(response.text)

    chapters = []
    items = rss_data["rss"]["channel"].get("item", [])
    if not isinstance(items, list):
        items = [items]

    vprint(f"\nDEBUG: Found {len(items)} items in RSS feed")

    for item in items:
        title = item["title"]
        url = item["link"]
        vprint(f"DEBUG: Processing RSS item: {title} -> {url}")

        # Extract chapter number from title
        # Handle both "Chapter X" and "Days X" formats
        chapter_num = None

        # For titles like "Sakamoto Days Days 196"
        if "Days" in title:
            parts = title.split("Days")
            if len(parts) > 1:
                num_part = parts[-1].strip()  # Get the last number after "Days"
                if num_part.replace(".", "").isdigit():
                    chapter_num = num_part
        # For titles with "Chapter X" format
        elif "Chapter" in title:
            parts = title.split("Chapter")
            if len(parts) > 1:
                num_part = parts[-1].strip()
                if num_part.replace(".", "").isdigit():
                    chapter_num = num_part

        if chapter_num:
            chapters.append({"chapter": chapter_num, "url": url})
            vprint(f"DEBUG: Added chapter {chapter_num} from {url}")

    # Sort chapters numerically
    sorted_chapters = sorted(
        chapters,
        key=lambda x: float(x["chapter"])
        if x["chapter"].replace(".", "").isdigit()
        else 0,
    )
    vprint(f"\nDEBUG: Final chapter count: {len(sorted_chapters)}")
    return sorted_chapters


def extract_chapter_images(chapter_url, chapter_num):
    """Extract image links from a chapter page using Playwright"""
    image_links = []

    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        context = browser.new_context(viewport={"width": 1920, "height": 1080})
        page = context.new_page()

        try:
            # Navigate to chapter URL and wait for content to load
            vprint(f"\nDEBUG: Fetching chapter from URL: {chapter_url}")
            page.goto(chapter_url, wait_until="networkidle", timeout=60000)

            # Wait for images to load
            page.wait_for_selector("img", state="visible", timeout=30000)

            # Extract all image URLs
            images = page.query_selector_all("img")
            for img in images:
                img_url = img.get_attribute("src")
                if img_url and img_url.lower().endswith(".png"):
                    # Skip brand.png
                    if not img_url.endswith("/static/images/brand.png"):
                        image_links.append(img_url)
                        vprint(f"Found image: {img_url}")

            vprint(f"\nFound {len(image_links)} images in Chapter {chapter_num}")
            return image_links

        except Exception as e:
            print(f"Error extracting chapter images: {str(e)}")
            return []
        finally:
            context.close()
            browser.close()


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Manga downloader script")
    parser.add_argument(
        "-o",
        "--output",
        type=str,
        default="manga_downloads",
        help="Output directory for downloaded manga (default: manga_downloads)",
    )
    parser.add_argument(
        "-z",
        "--zip",
        action="store_true",
        help="Create zip archives for chapters in volume format",
    )
    parser.add_argument(
        "-v", "--verbose", action="store_true", help="Enable verbose debug output"
    )
    parser.add_argument(
        "--no-skip", action="store_true", help="Do not skip existing zip files"
    )
    parser.add_argument(
        "--chapter-filter",
        type=float,
        metavar="N",
        help="Download only chapters ABOVE the given number (Example: --chapter-filter 46 will download chapters 47 and up)",
    )
    parser.add_argument(
        "-l",
        "--latest",
        action="store_true",
        help="Continue from latest downloaded chapter (automatically detects highest chapter number from existing zip files)",
    )
    parser.add_argument(
        "-t",
        "--title",
        type=str,
        default="wistoria",
        help="Manga title to search for (default: wistoria)",
    )
    parser.add_argument(
        "-b",
        "--bulk",
        type=str,
        help="Path to text file containing manga titles (one per line)",
    )
    args = parser.parse_args()

    def process_manga_title(title, filter_chapter=None):
        """Process a single manga title"""
        # Replace hyphens with spaces in the manga title
        search_title = title.replace("-", " ")
        print(f"\nProcessing manga: {search_title}")

        manga_url = search_manga(search_title)
        if not manga_url:
            print(f"Could not find manga: {search_title}")
            return

        # Create manga-specific directory using the slug
        manga_slug = get_manga_slug(manga_url)
        manga_dir = os.path.join(args.output, manga_slug) if manga_slug else args.output

        if args.latest and os.path.exists(manga_dir):
            # Find zip files only in this manga's directory
            zip_files = [
                f
                for f in os.listdir(manga_dir)
                if f.startswith("vol_") and f.endswith(".zip")
            ]

            if zip_files:
                chapter_numbers = []
                for zip_file in zip_files:
                    try:
                        chapter_num = float(zip_file[4:-4])
                        chapter_numbers.append(chapter_num)
                    except ValueError:
                        continue

                if chapter_numbers:
                    filter_chapter = max(chapter_numbers)
                    print(f"Found highest chapter {filter_chapter} in {manga_dir}")
                    print(f"Will only download chapters above chapter {filter_chapter}")

        # Create manga-specific directory
        os.makedirs(manga_dir, exist_ok=True)

        # Pre-scan existing zip files
        existing_zips = set()
        if args.zip and not args.no_skip and os.path.exists(manga_dir):
            existing_zips = {f for f in os.listdir(manga_dir) if f.endswith(".zip")}
            vprint(f"Found {len(existing_zips)} existing zip files")

        # Rest of your existing manga processing code...
        has_existing_chapters = False
        if os.path.exists(manga_dir):
            dir_contents = os.listdir(manga_dir)
            has_existing_chapters = any(
                f.startswith(("vol_", "chapter_")) for f in dir_contents
            )

        if has_existing_chapters:
            print("Existing manga chapters found, using RSS feed for updates...")
            rss_url = get_rss_url(manga_url)
            print(f"RSS URL: {rss_url}")
            chapters = get_chapter_links(rss_url)
        else:
            print("New manga or empty folder, fetching full chapter list...")
            chapter_list_url = get_chapter_list_url(manga_url)
            print(f"Chapter list URL: {chapter_list_url}")
            chapters = get_chapters_from_list(chapter_list_url)

        # Process chapters...
        if chapters:
            min_chapter = min(
                float(c["chapter"])
                for c in chapters
                if c["chapter"].replace(".", "").isdigit()
            )
            max_chapter = max(
                float(c["chapter"])
                for c in chapters
                if c["chapter"].replace(".", "").isdigit()
            )
            print(f"\nFound chapters from {min_chapter} to {max_chapter}")
            if filter_chapter:
                print(
                    f"Will download chapters from {filter_chapter + 1} to {max_chapter}"
                )

            print("\nChapters found:")
            for chapter in chapters:
                print(f"Chapter {chapter['chapter']} - {chapter['url']}")
                try:
                    chapter_num = float(chapter["chapter"])

                    if filter_chapter is not None:
                        if chapter_num <= filter_chapter:
                            vprint(
                                f"Skipping chapter {chapter_num} (not higher than {filter_chapter})"
                            )
                            continue
                        else:
                            print(f"Processing chapter {chapter_num}")
                except ValueError:
                    print(
                        f"Warning: Could not parse chapter number from {chapter['chapter']}"
                    )
                    continue

                # Handle chapter versions
                base_chapter = int(float(chapter_num))
                version = None
                if "." in str(chapter_num):
                    version = int(str(chapter_num).split(".")[1])

                vol_name = f"vol_{base_chapter:03d}"
                if version:
                    vol_name = f"{vol_name}-{version}"
                zip_filename = f"{vol_name}.zip"

                if args.zip and not args.no_skip and zip_filename in existing_zips:
                    vprint(f"Skipping existing zip archive: {zip_filename}")
                    continue

                image_links = extract_chapter_images(chapter["url"], chapter["chapter"])
                vprint(
                    f"Found {len(image_links)} images in Chapter {chapter['chapter']}"
                )

                chapter_dir = os.path.join(manga_dir, f"chapter_{chapter['chapter']}")
                downloaded_files = download_images_parallel(image_links, chapter_dir)
                vprint(f"Successfully downloaded {len(downloaded_files)} images")

                if args.zip:
                    zip_base = os.path.join(manga_dir, vol_name)
                    zip_path = f"{zip_base}.zip"
                    print(f"Creating zip archive: {zip_path}")
                    shutil.make_archive(zip_base, "zip", chapter_dir)
                    shutil.rmtree(chapter_dir)

    filter_chapter = args.chapter_filter

    if args.bulk:
        try:
            with open(args.bulk, "r", encoding="utf-8") as f:
                manga_titles = [line.strip() for line in f if line.strip()]
            print(f"Found {len(manga_titles)} manga titles in {args.bulk}")
            for title in manga_titles:
                process_manga_title(title, filter_chapter)
        except FileNotFoundError:
            print(f"Error: Could not find bulk file: {args.bulk}")
            exit(1)
    else:
        process_manga_title(args.title, filter_chapter)

    def vprint(*print_args, **kwargs):
        """Print only if verbose mode is enabled"""
        if args.verbose:
            print(*print_args, **kwargs)
