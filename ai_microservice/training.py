import sqlite3
import pandas as pd
from sklearn.ensemble import RandomForestRegressor
from sklearn.model_selection import train_test_split
import pickle
import os


town_names = [
    "ANG MO KIO", "BEDOK", "BISHAN", "BUKIT BATOK", "BUKIT MERAH", "BUKIT PANJANG",
    "BUKIT TIMAH", "CENTRAL AREA", "CHOA CHU KANG", "CLEMENTI", "GEYLANG",
    "HOUGANG", "JURONG EAST", "JURONG WEST", "KALLANG/WHAMPOA", "LIM CHU KANG",
    "MARINE PARADE", "PASIR RIS", "PUNGGOL", "QUEENSTOWN", "SEMBAWANG",
    "SENGKANG", "SERANGOON", "TAMPINES", "TOA PAYOH", "WOODLANDS", "YISHUN"
]

conn = sqlite3.connect("../backend/hdb.db")
query = """
SELECT time, flat_type, resale_price
FROM hdb_records
WHERE town = ?
AND flat_type = "4 ROOM"
"""

town_data = {}
for town in town_names:
    df = pd.read_sql_query(query, conn, params=(town,))
    town_data[town] = df
conn.close()

os.makedirs("models", exist_ok=True)

for town, data in town_data.items():   
    if data.empty:
        print(f"No data available for town: {town}")
        continue

    data['time'] = pd.to_datetime(data['time'])
    data['time'] = data['time'].view('int64') // 10**9 
    X = data[['time']]
    y = data['resale_price']

    X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)

    model = RandomForestRegressor(n_estimators=100, random_state=42)
    model.fit(X_train, y_train)

    model_path = os.path.join("models", f"{town.replace('/', '_')}_model.pkl")
    with open(model_path, "wb") as file:
        pickle.dump(model, file)

    print(f"Model saved for town: {town}")