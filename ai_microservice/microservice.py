from flask import Flask, request, jsonify, send_file
import pandas as pd
import joblib
import datetime
import matplotlib
matplotlib.use("Agg")
import matplotlib.pyplot as plt
import io

app = Flask(__name__)

flat_type_mapping = {
    '2 ROOM': 1,
    '3 ROOM': 2,
    '4 ROOM': 3,
    '5 ROOM': 4,
    'EXECUTIVE': 5
}

@app.route("/predict", methods=["POST"])
def predict():
    town = request.args.get("town", default="BEDOK", type=str)
    flat_type = request.args.get("flat_type", default="4 ROOM", type=str)

    model_path = f"models/{town}_model.pkl"
    try:
        model = joblib.load(model_path)
    except FileNotFoundError:
        return jsonify({"error": f"Model for town '{town}' not found."}), 404
    
    if flat_type not in flat_type_mapping:
        return jsonify({"error": f"Invalid flat_type '{flat_type}'. Valid options: {list(flat_type_mapping.keys())}"}), 400
    
    flat_type_num = flat_type_mapping[flat_type]

    start_date = datetime.date(2020, 1, 1)
    end_date = datetime.date(2021, 1, 1)
    date_range = pd.date_range(start_date, end_date, freq="D")

    input_data = pd.DataFrame({
        "time": date_range,
        "flat_type_num": flat_type_num,
    })

    input_data["time"] = input_data["time"].map(pd.Timestamp.timestamp)

    predictions = model.predict(input_data)

    plt.figure(figsize=(12, 6))
    plt.plot(date_range, predictions, label=f"Predicted Prices ({flat_type})", color='blue', marker='o', markersize=2)
    plt.title(f"Predicted Prices for {town} ({flat_type})")
    plt.xlabel("Date")
    plt.ylabel("Resale Price")
    plt.legend()
    plt.grid()

    img = io.BytesIO()
    plt.savefig(img, format='png')
    img.seek(0)
    plt.close()

    response_data = {
        "town": town,
        "flat_type": flat_type,
        "predictions": [
            {"date": date.strftime("%Y-%m-%d"), "predicted_price": float(price)}
            for date, price in zip(date_range, predictions)
        ]
    }

    return send_file(img, mimetype='image/png')

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5433)