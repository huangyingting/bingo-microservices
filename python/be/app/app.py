# curl -XPOST http://localhost:5000/v1/extract \
#   --header 'content-type: application/json' \
#   --data '{"url": "https://www.microsoft.com"}'

from flask import Flask, jsonify, request
from newspaper import Article
from opentelemetry.instrumentation.flask import FlaskInstrumentor
from opentelemetry import trace
from opentelemetry.exporter.jaeger.thrift import JaegerExporter
from opentelemetry.sdk.resources import SERVICE_NAME, Resource
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.instrumentation.requests import RequestsInstrumentor
import os
from prometheus_flask_exporter.multiprocess import GunicornInternalPrometheusMetrics

trace.set_tracer_provider(
TracerProvider(
        resource=Resource.create({SERVICE_NAME: "be"})
    )
)
tracer = trace.get_tracer(__name__)
jaeger_exporter = JaegerExporter(
    collector_endpoint=os.getenv('BE_JAEGER_ADDR', 'http://localhost:14268/api/traces'),
)
span_processor = BatchSpanProcessor(jaeger_exporter)
trace.get_tracer_provider().add_span_processor(span_processor)

app = Flask(__name__)
metrics = GunicornInternalPrometheusMetrics(app, defaults_prefix="bingo_be")
FlaskInstrumentor().instrument_app(app)
RequestsInstrumentor().instrument()

@app.route('/healthz')
@metrics.do_not_track()
def healthz():
  response = {
      "status": "OK"
  }  
  return jsonify(response), 200


@app.route('/readyz')
@metrics.do_not_track()
def readyz():
  response = {
      "status": "OK"
  }  
  return jsonify(response), 200

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


metrics.register_default(
    metrics.counter(
        'by_path_counter', 'Request count by request paths',
        labels={'path': lambda: request.path}
    )
)

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8002, debug=True)
