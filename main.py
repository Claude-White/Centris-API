import json

with open("data.json", "r", encoding="utf-8") as file:
    data = json.load(file)

mls_numbers = [item["mls"] for item in data]
duplicates = {mls for mls in mls_numbers if mls_numbers.count(mls) > 1}

print("Duplicate MLS Numbers:", duplicates)
