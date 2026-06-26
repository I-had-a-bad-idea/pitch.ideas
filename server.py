from flask import Flask, render_template
import db

app = Flask(__name__)

@app.route("/")
def home():
    return render_template("index.html")

if __name__ == "__main__":
    db.init_db()
    app.run(host="localhost", port=4000)