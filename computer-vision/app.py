import re
from flask import Flask, request, jsonify, send_file
import os
import cv2
import numpy as np
from sklearn.cluster import KMeans
from datetime import datetime

app = Flask(__name__)

IMAGE_PATH = '/app/images'

def preprocess_image(image_path):
    # Load the image
    image = cv2.imread(image_path)
    if image is None:
        raise ValueError("Image not found or unable to load.")
    
    # Convert image to RGB (OpenCV uses BGR by default)
    image_rgb = cv2.cvtColor(image, cv2.COLOR_BGR2RGB)
    return image_rgb

def kmeans_clustering(image_rgb, n_clusters=2):
    # Reshape image for clustering
    pixels = image_rgb.reshape(-1, 3)
    
    # Perform K-means clustering
    kmeans = KMeans(n_clusters=n_clusters, random_state=42)
    kmeans.fit(pixels)
    
    # Get cluster labels and reshape to original image shape
    labels = kmeans.labels_.reshape(image_rgb.shape[:2])
    return labels

def enhance_mask(image_rgb, n_clusters=2):
    # Perform K-means clustering
    labels = kmeans_clustering(image_rgb, n_clusters)
    
    # Assume that the cluster with the highest mean green value represents the forest
    cluster_means = np.array([np.mean(image_rgb[labels == i], axis=0) for i in range(n_clusters)])
    green_cluster = np.argmax(cluster_means[:, 1])  # Assuming green channel is at index 1
    
    # Create mask for the forest cluster
    mask = (labels == green_cluster).astype(np.uint8) * 255
    
    # Apply morphological operations to clean up the mask
    kernel = np.ones((5, 5), np.uint8)
    mask_cleaned = cv2.morphologyEx(mask, cv2.MORPH_CLOSE, kernel)
    mask_cleaned = cv2.morphologyEx(mask_cleaned, cv2.MORPH_OPEN, kernel)

    return mask_cleaned

def calculate_forest_coverage(mask):
    # Count non-zero (white) pixels which represent the forest
    forest_pixels = np.sum(mask == 0)  # Pixels where mask is white (forest)
    
    # Count zero (black) pixels which represent non-forest
    non_forest_pixels = np.sum(mask == 255)  # Pixels where mask is black (non-forest)
    
    total_pixels = mask.size
    
    forest_percentage = (forest_pixels / total_pixels) * 100
    non_forest_percentage = (non_forest_pixels / total_pixels) * 100

    return forest_percentage

def extract_area_number(filename):
    # Define the regular expression pattern
    pattern = r'area_(\d+)_\d+\.png'
    
    # Use regular expression to search for the pattern in the filename
    match = re.search(pattern, filename)
    
    if match:
        # Extract the number from the match object
        number = match.group(1)
        return number
    else:
        # Return None or an appropriate message if no match is found
        return None

@app.route('/calculate-deforestation/<path:image_name>', methods=['GET'])
def calculate_deforestation(image_name):
    image_path = os.path.join(IMAGE_PATH, image_name)
    
    try:
        image_rgb = preprocess_image(image_path)

        # Perform segmentation
        mask = enhance_mask(image_rgb)

        # Calculate percentage of forest and non-forest areas
        forest_coverage = calculate_forest_coverage(mask)

        # Create a green image for the mask
        green_overlay = np.zeros_like(image_rgb)
        
        # Use mask to highlight forest areas
        # We want to highlight the areas where the mask is black (forest) in green
        green_overlay[mask == 0] = [255, 255, 0]  # Green color for forest

        # Blend the overlay with the original image
        alpha = 0.33  # Transparency factor
        overlay_image = cv2.addWeighted(image_rgb, 1 - alpha, green_overlay, alpha, 0)
        
        # Save the mask image
        area_id = extract_area_number(image_name)  # Dummy value; replace with actual area ID logic
        timestamp = datetime.now().strftime('%Y%m%d%H%M%S')
        masked_image_path = os.path.join(IMAGE_PATH, f"masked_area_{area_id}_{timestamp}.png")
        cv2.imwrite(masked_image_path, overlay_image)

        return jsonify({
            'forest_coverage': forest_coverage,
            'masked_image_path': masked_image_path
        })
    except Exception as e:
        return jsonify({'error': str(e)}), 400

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
