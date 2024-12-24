import math
import json

# Input data
big_square = {
    "SouthWest": {"Lat": 44.980811, "Lng": -79.670110},
    "NorthEast": {"Lat": 51.998494, "Lng": -57.107759}
}

chunk_size = {"Lat": 0.6141262171, "Lng": 0.835647583}

# Initialize variables
chunks = []
current_lat = big_square["SouthWest"]["Lat"]
south_west_lat_limit = big_square["NorthEast"]["Lat"]
south_west_lng_limit = big_square["NorthEast"]["Lng"]

# Generate chunks
while current_lat < south_west_lat_limit:
    current_lng = big_square["SouthWest"]["Lng"]
    while current_lng < south_west_lng_limit:
        chunk = {
            "zoomLevel": 11,
            "mapBounds": {
                "SouthWest": {"Lat": current_lat, "Lng": current_lng},
                "NorthEast": {
                    "Lat": min(current_lat + chunk_size["Lat"], south_west_lat_limit),
                    "Lng": min(current_lng + chunk_size["Lng"], south_west_lng_limit)
                }
            }
        }
        chunks.append(chunk)
        current_lng += chunk_size["Lng"]
    current_lat += chunk_size["Lat"]

# Output result
print(json.dumps(chunks, indent=2))
