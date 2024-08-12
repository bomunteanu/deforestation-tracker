import math
import requests
from PIL import Image
from io import BytesIO

API_KEY = 'egZCFjVhe1y2Mp3AmEonfYsyLqx0mq0c'

def lat_lon_to_tile(lat, lon, zoom):
    lat_rad = math.radians(lat)
    n = 2.0 ** zoom
    x = int((lon + 180.0) / 360.0 * n)
    y = int((1.0 - math.log(math.tan(lat_rad) + 1.0 / math.cos(lat_rad)) / math.pi) / 2.0 * n)
    return x, y

def get_satellite_image_tile(z, x, y):
    API_ENDPOINT = f"https://api.tomtom.com/map/1/tile/sat/main/{z}/{x}/{y}.jpg?key={API_KEY}"
    response = requests.get(API_ENDPOINT)
    if response.status_code == 200:
        print("OK")
        return response.content
    else:
        print(f"Failed to retrieve tile ({z}/{x}/{y}). Status code: {response.status_code}")
        return None

def download_tiles(lat_min, lon_min, lat_max, lon_max, zoom):
    x_min, y_min = lat_lon_to_tile(lat_max, lon_min, zoom)  # Bottom-left
    x_max, y_max = lat_lon_to_tile(lat_min, lon_max, zoom)  # Top-right
    
    print(f"Tile range: x_min={x_min}, y_max={y_max}, x_max={x_max}, y_min={y_min}")
    
    tiles = []
    tmp = []
    for y in range(y_min, y_max + 1):
        row_tiles = []
        row_tmp = []
        for x in range(x_min, x_max + 1):  # x varies from left to right
            print(f"Downloading tile ({zoom}, {x}, {y})")
            tile_image = get_satellite_image_tile(zoom, x, y)
            if tile_image:
                row_tiles.append(tile_image)
                row_tmp.append((x, y))
            else:
                print(f"Tile ({zoom}, {x}, {y}) failed to download.")
        tiles.append(row_tiles)
        tmp.append(row_tmp)
    
    if not any(row for row in tiles):  # Check if all rows are empty
        print("No tiles were downloaded.")
    
    return tiles, x_min, y_max, x_max, y_min, tmp


def stitch_tiles(tiles, x_min, y_max, x_max, y_min, tmp, tile_size=256):

    if not tiles:
        print("No tiles to stitch.")
        return None
    
    num_rows = len(tiles)
    num_cols = len(tiles[0])
    
    if num_rows == 0 or num_cols == 0:
        print("Error: No tiles found.")
        return None
    
    width = num_cols * tile_size
    height = num_rows * tile_size
    
    print(f"Creating stitched image with dimensions: {width}x{height}")
    
    stitched_image = Image.new('RGB', (width, height))
    
    for row_idx, row_tiles in enumerate(tiles):
        for col_idx, tile_image in enumerate(row_tiles):
            print(tmp[row_idx][col_idx])
            print("_____________________________")
            try:
                image = Image.open(BytesIO(tile_image))
                if image.size != (tile_size, tile_size):
                    print(f"Warning: Tile {row_idx}, {col_idx} has size {image.size} instead of expected {tile_size}x{tile_size}")
                x_pos = col_idx * tile_size
                y_pos = (row_idx) * tile_size  # Correct y position for top-down stitching
                print(x_pos, y_pos)
                if x_pos + tile_size > width or y_pos + tile_size > height:
                    print(f"Error: Tile {row_idx}, {col_idx} extends outside the bounds of the final image.")
                    continue
                stitched_image.paste(image, (x_pos, y_pos))
            except Exception as e:
                print(f"Error processing tile {row_idx}, {col_idx}: {e}")
    
    return stitched_image


def main():
    zoom = 15
    lat_min = 45.369610  # Bottom left latitude
    lon_min = 25.693347  # Bottom left longitude
    lat_max = 45.389264
    lon_max = 25.763041  # Top right longitude
    
    tiles, x_min, y_max, x_max, y_min, tmp = download_tiles(lat_min, lon_min, lat_max, lon_max, zoom)
    
    print(f"Tiles downloaded: {len(tiles)} rows, {len(tiles[0]) if tiles else 0} columns")
    
    stitched_image = stitch_tiles(tiles, x_min, y_max, x_max, y_min, tmp)
    
    if stitched_image:
        try:
            stitched_image.save('stitched_image.png')
            print("Stitched image saved as 'stitched_image.png'")
        except Exception as e:
            print(f"Error saving stitched image: {e}")
    else:
        print("Failed to create stitched image.")

if __name__ == "__main__":
    main()
