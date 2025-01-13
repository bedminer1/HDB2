import sqlite3
import pandas as pd
import numpy as np
from sklearn.model_selection import train_test_split
from tensorflow.keras.models import Sequential
from tensorflow.keras.layers import LSTM, Dense
from sklearn.preprocessing import MinMaxScaler
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

    # Preprocess time and flat_type
    data['time'] = pd.to_datetime(data['time'])
    data['time'] = data['time'].astype('int64') // 10**9
    data['flat_type_num'] = data['flat_type'].map(flat_type_mapping)
    
    # Normalize the data
    scaler = MinMaxScaler()
    data[['time', 'flat_type_num', 'resale_price']] = scaler.fit_transform(data[['time', 'flat_type_num', 'resale_price']])
    
    # Prepare sequences for LSTM
    sequence_length = 30  # Use the past 30 days to predict the next day
    X, y = [], []
    for i in range(len(data) - sequence_length):
        X.append(data[['time', 'flat_type_num']].iloc[i:i + sequence_length].values)
        y.append(data['resale_price'].iloc[i + sequence_length])
    X, y = np.array(X), np.array(y)

    # Train-test split
    X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)

    # Define the LSTM model
    model = Sequential([
        LSTM(50, activation='relu', input_shape=(X_train.shape[1], X_train.shape[2])),
        Dense(1)
    ])
    model.compile(optimizer='adam', loss='mse')

    # Train the LSTM model
    model.fit(X_train, y_train, epochs=20, batch_size=32, verbose=1, validation_data=(X_test, y_test))

    # Save the model
    model_path = os.path.join("models", f"{town.replace('/', '_')}_model.h5")
    model.save(model_path)

    # Save the scaler for normalization
    scaler_path = os.path.join("models", f"{town.replace('/', '_')}_scaler.pkl")
    with open(scaler_path, "wb") as file:
        pickle.dump(scaler, file)

    print(f"LSTM model and scaler saved for town: {town}")