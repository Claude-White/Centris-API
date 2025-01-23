import json

def load_json(filename):
    with open(filename, "r", encoding="utf-8") as file:
        return json.load(file)

def flatten(data):
    """Recursively flattens a list of lists into a single list of dictionaries."""
    flat_list = []
    for item in data:
        if isinstance(item, list):
            flat_list.extend(flatten(item))  # Recursively flatten nested lists
        elif isinstance(item, dict):  # Only append valid dictionaries
            flat_list.append(item)
    return flat_list

def validate_property_ids(properties_file, *reference_files):
    properties = load_json(properties_file)
    reference_data = [flatten(load_json(f)) for f in reference_files]
    
    property_ids = {prop["mls"] for prop in properties if isinstance(prop, dict)}

    missing_references = []
    
    for ref_file, data in zip(reference_files, reference_data):
        for item in data:
            if isinstance(item, dict) and "property_id" in item:  # Ensure item is valid
                if item["property_id"] not in property_ids:
                    missing_references.append((ref_file, item["property_id"]))
    
    if missing_references:
        print("Missing property IDs found:")
        for ref_file, prop_id in missing_references:
            print(f"File: {ref_file}, Missing property_id: {prop_id}")
    else:
        print("All property references are valid.")

if __name__ == "__main__":
    validate_property_ids("properties.json", "propertiesExpenses.json")
