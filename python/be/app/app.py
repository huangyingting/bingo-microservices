# curl -XPOST http://localhost:5000/v1/extract \
#   --header 'content-type: application/json' \
#   --data '{"url": "https://www.microsoft.com"}'

from flask import Flask, jsonify, request
from newspaper import Article

app = Flask(__name__)


@app.route('/v1/extract', methods=['POST'])
def extract_html():
    if not request.json or not 'url' in request.json:
        return jsonify({"code": 400, "reason": "INVALID_REQUEST_PAYLOAD", "message": "invalid request payload"}), 400
    article = Article(url=request.json['url'], language='en')
    article.download()
    article.parse()
    article.nlp()
    response = {
        "title": str(article.title),
        "authors": list(article.authors),
        "published_date": str(article.publish_date),
        "videos": list(article.movies),
        "keywords": list(article.keywords),
        "tags": list(article.tags),
        "meta_keywords": list(article.meta_keywords),
        "summary": str(article.summary)
    }

    return jsonify(response), 200


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8002, debug=True)
