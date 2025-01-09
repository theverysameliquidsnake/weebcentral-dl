import requests
import os
from urllib.parse import urlparse
from concurrent.futures import ThreadPoolExecutor, as_completed


def download_single_image(url, save_dir, headers=None):
    """
    Download an image from a URL and save it to a specified directory.

    Args:
        url (str): URL of the image
        save_dir (str): Directory to save the image (default: 'downloads')
        headers (dict): Request headers (default: None)

    Returns:
        str: Path to the saved image file if successful, None otherwise
    """
    try:
        # Create default headers if none provided
        if headers is None:
            headers = {
                "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0"
            }

        # Send GET request
        response = requests.get(url, headers=headers, stream=True)
        response.raise_for_status()  # Raise exception for bad status codes

        # Create save directory if it doesn't exist
        os.makedirs(save_dir, exist_ok=True)

        # Extract filename from URL
        filename = os.path.basename(urlparse(url).path)
        if not filename:
            filename = "image.png"  # Default filename if none found in URL

        # Full path for saving the file
        save_path = os.path.join(save_dir, filename)

        # Save the image
        with open(save_path, "wb") as f:
            for chunk in response.iter_content(chunk_size=8192):
                if chunk:
                    f.write(chunk)

        # Successfully downloaded without printing
        return save_path

    except requests.exceptions.RequestException as e:
        print(f"Error downloading image: {e}")
        return None


def download_images_parallel(urls, save_dir, max_workers=4):
    """
    Download multiple images in parallel using ThreadPoolExecutor.

    Args:
        urls (list): List of image URLs to download
        save_dir (str): Directory to save the images
        max_workers (int): Maximum number of concurrent downloads

    Returns:
        list: List of successfully downloaded image paths
    """
    downloaded_files = []

    with ThreadPoolExecutor(max_workers=max_workers) as executor:
        future_to_url = {
            executor.submit(download_single_image, url, save_dir): url for url in urls
        }

        for future in as_completed(future_to_url):
            url = future_to_url[future]
            try:
                result = future.result()
                if result:
                    downloaded_files.append(result)
                    print(f"Successfully downloaded: {url}")
            except Exception as e:
                print(f"Failed to download {url}: {str(e)}")

    return downloaded_files
