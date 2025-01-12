import sqlite3
import pandas as pd
import joblib
from sklearn.metrics import mean_absolute_error, root_mean_squared_error
import matplotlib.pyplot as plt

# Define a simple mapping for flat_type
flat_type_mapping = {
    "2 ROOM": 1,
    "3 ROOM": 2,
    "4 ROOM": 3,
    "5 ROOM": 4,
    "EXECUTIVE": 5
}

town_to_predict = "ANG MO KIO"

# Fetch actual data for prediction
conn = sqlite3.connect("../backend/hdb.db")
query = """
SELECT time, flat_type, resale_price
FROM hdb_records
WHERE town = ?
"""
actual_data = pd.read_sql_query(query, conn, params=(town_to_predict,))
conn.close()

# Preprocess the data
actual_data['time'] = pd.to_datetime(actual_data['time'])
actual_data['time'] = actual_data['time'].astype('int64') // 10**9  # Convert time to seconds since epoch
actual_data['flat_type'] = actual_data['flat_type'].map(flat_type_mapping)  # Map flat_type to numeric

actual_data.sort_values(by='time', inplace=True)

# Load the trained model
model_path = f"models/{town_to_predict}_model.pkl"
model = joblib.load(model_path)

# Prepare features for prediction
X_test = actual_data[['time', 'flat_type']].values  # Include both time and flat_type
y_test = actual_data['resale_price']

# Make predictions
predictions = model.predict(X_test)
actual_data['predicted_price'] = predictions

# Calculate error metrics
mae = mean_absolute_error(y_test, predictions)
rmse = root_mean_squared_error(y_test, predictions)

print(f"Mean Absolute Error (MAE): {mae:.2f}")
print(f"Root Mean Squared Error (RMSE): {rmse:.2f}")

# Plot the results
plt.figure(figsize=(12, 6))
plt.plot(actual_data['time'], actual_data['resale_price'], label="Actual Price", marker='o')
plt.plot(actual_data['time'], actual_data['predicted_price'], label="Predicted Price", marker='x')
plt.title(f"Actual vs Predicted Prices for {town_to_predict}")
plt.xlabel("Time")
plt.ylabel("Resale Price")
plt.legend()
plt.grid()
plt.show()