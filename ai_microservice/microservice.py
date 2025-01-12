from flask import Flask, request, jsonify

app = Flask(__name__)

@app.route("/predict", methods=["POST"])
def predict():
    town = request.args.get("town", default="BEDOK", type=str)
    flat_type = request.args.get("flat_type", default="4 ROOM", type=str)



    data = {}
    data["message"] = f"the town to be analyzed is {town}, {flat_type}"

    return jsonify(data)

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5433)