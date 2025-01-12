import sqlite3
import pandas as pd
import joblib
from sklearn.metrics import mean_absolute_error, root_mean_squared_error
import matplotlib.pyplot as plt

town_to_predict = "ANG MO KIO"

conn = sqlite3.connect("../backend/hdb.db")
query = """
SELECT time, resale_price
FROM hdb_records
WHERE town = ?
AND flat_type = "4 ROOM"
"""
actual_data = pd.read_sql_query(query, conn, params=(town_to_predict,))
conn.close()
actual_data['time'] = pd.to_datetime(actual_data['time'])
actual_data.sort_values(by='time', inplace=True)

# Load a model
model_path = f"models/{town_to_predict}_model.pkl"
model = joblib.load(model_path)

# Predict for new data
X_test = actual_data['time'].map(pd.Timestamp.timestamp).values.reshape(-1, 1)

predictions = model.predict(X_test)
actual_data['predicted_price'] = predictions

mae = mean_absolute_error(actual_data['resale_price'], actual_data['predicted_price'])
rmse = root_mean_squared_error(actual_data['resale_price'], actual_data['predicted_price'])
print(f"Mean Absolute Error (MAE): {mae:.2f}")
print(f"Root Mean Squared Error (RMSE): {rmse:.2f}")

# plt.figure(figsize=(12, 6))
# plt.plot(actual_data['time'], actual_data['resale_price'], label="Actual Price", marker='o')
# plt.plot(actual_data['time'], actual_data['predicted_price'], label="Predicted Price", marker='x')
# plt.title(f"Actual vs Predicted Prices for {town_to_predict} (4-Room)")
# plt.xlabel("Time")
# plt.ylabel("Resale Price")
# plt.legend()
# plt.grid()
# plt.show()