from flask import Flask,jsonify
import json
app = Flask(__name__)

with open("responses.txt", "r") as f:
    contents = json.load(f)

@app.route("/status")
def hello():
    return jsonify(contents.pop())
