import sqlite3
import pandas as pd
from sklearn.ensemble import RandomForestRegressor
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import OneHotEncoder
import pickle
import os

# List of towns
town_names = [
    "ANG MO KIO", "BEDOK", "BISHAN", "BUKIT BATOK", "BUKIT MERAH", "BUKIT PANJANG",
    "BUKIT TIMAH", "CENTRAL AREA", "CHOA CHU KANG", "CLEMENTI", "GEYLANG",
    "HOUGANG", "JURONG EAST", "JURONG WEST", "KALLANG/WHAMPOA", "LIM CHU KANG",
    "MARINE PARADE", "PASIR RIS", "PUNGGOL", "QUEENSTOWN", "SEMBAWANG",
    "SENGKANG", "SERANGOON", "TAMPINES", "TOA PAYOH", "WOODLANDS", "YISHUN"
]

flat_type_mapping = {
    '2 ROOM': 1,
    '3 ROOM': 2,
    '4 ROOM': 3,
    '5 ROOM': 4,
    'EXECUTIVE': 5
}

# Connect to the SQLite database
conn = sqlite3.connect("../backend/hdb.db")
query = """
SELECT time, flat_type, resale_price
FROM hdb_records
WHERE town = ?
"""

town_data = {}
for town in town_names:
    df = pd.read_sql_query(query, conn, params=(town,))
    town_data[town] = df
conn.close()

# Create models folder if it doesn't exist
os.makedirs("models", exist_ok=True)

# Process data and train models town by town
for town, data in town_data.items():   
    if data.empty:
        print(f"No data available for town: {town}")
        continue

    # Preprocess time
    data['time'] = pd.to_datetime(data['time'])
    data['time'] = data['time'].astype('int64') // 10**9
    data['flat_type_num'] = data['flat_type'].map(flat_type_mapping)
    
    X = data[['time', 'flat_type_num']]
    y = data['resale_price']

    # Train-test split
    X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)

    # Train the RandomForestRegressor
    model = RandomForestRegressor(n_estimators=100, random_state=42)
    model.fit(X_train, y_train)

    # Save the model and encoder
    model_path = os.path.join("models", f"{town.replace('/', '_')}_model.pkl")
    encoder_path = os.path.join("models", f"{town.replace('/', '_')}_encoder.pkl")

    with open(model_path, "wb") as file:
        pickle.dump(model, file)

    print(f"Model saved for town: {town}")
  