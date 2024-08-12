export function calculateRectangleArea(
    topRightLat: number,
    topRightLon: number,
    bottomLeftLat: number,
    bottomLeftLon: number
): number {
    const R = 6371; // Radius of the Earth in kilometers

    // Convert degrees to radians
    const toRadians = (degrees: number) => degrees * Math.PI / 180;

    // Calculate the distance between two latitudes and two longitudes
    const latDistance = toRadians(topRightLat - bottomLeftLat);
    const lonDistance = toRadians(topRightLon - bottomLeftLon);

    // Calculate the average latitude for the longitudinal distance calculation
    const avgLat = toRadians((topRightLat + bottomLeftLat) / 2);

    // Haversine formula for distance calculation
    const a = Math.sin(latDistance / 2) ** 2 +
              Math.cos(toRadians(bottomLeftLat)) * Math.cos(toRadians(topRightLat)) *
              Math.sin(lonDistance / 2) ** 2;
    const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
    
    // Distance in kilometers
    const latDistanceKm = R * c;

    // Adjust the longitudinal distance based on the average latitude
    const lonDistanceKm = R * Math.abs(lonDistance) * Math.cos(avgLat);

    // Calculate area
    const area = latDistanceKm * lonDistanceKm;

    return area;
}