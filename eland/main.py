import elasticsearch
from pathlib import Path
from eland.ml.pytorch import PyTorchModel
from eland.ml.pytorch.transformers import TransformerModel

# Load a Hugging Face transformers model directly from the model hub
tm = TransformerModel("sentence-transformers/msmarco-MiniLM-L-12-v3", "text_embedding")

# Export the model in a TorchScrpt representation which Elasticsearch uses
tmp_path = "models"
Path(tmp_path).mkdir(parents=True, exist_ok=True)
model_path, config, vocab_path = tm.save(tmp_path)

# model_path = "models/traced_pytorch_model.pt"
# vocab_path = "models/vocabulary.json"

# Import model into Elasticsearch
es = elasticsearch.Elasticsearch("https://elastic:elastic@localhost:9200", request_timeout=300, ca_certs=Path("/tmp/ca.crt"))  # 5 minute timeout
ptm = PyTorchModel(es, tm.elasticsearch_model_id())
ptm.import_model(model_path=model_path, config_path=None, vocab_path=vocab_path, config=config)
