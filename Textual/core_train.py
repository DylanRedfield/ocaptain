from rasa_core.agent import Agent
from rasa_core.interpreter import RasaNLUInterpreter
from rasa_core import utils
from rasa_core.train import train_dialogue_model
#from rasa_core.utils import AvailableEndpoints
from rasa_core.run import AvailableEndpoints
from rasa_core.constants import (
        DEFAULT_NLU_FALLBACK_THRESHOLD,
        DEFAULT_CORE_FALLBACK_THRESHOLD, DEFAULT_FALLBACK_ACTION)

interpreter = RasaNLUInterpreter("models/current/nlu")

#endpoint_conf = utils.read_endpoint_config("endpoints.yml", endpoint_type="action_endpoint")
#endpoints = AvailableEndpoints(nlg = None, nlu = None, action = endpoint_conf, model = None)

default_args = {
        "epochs" : 200,
        "validation_split": 0.1,
        "batch_size" : 20,
        "augmentation_factor": 50,
        "nlu_threshold": DEFAULT_NLU_FALLBACK_THRESHOLD,
        "core_threshold": DEFAULT_CORE_FALLBACK_THRESHOLD,
        "fallback_action_name": DEFAULT_FALLBACK_ACTION
        }

agent = train_dialogue_model("domain.yml", "stories.md", "models/dialogue", None, "endpoints.yml", 3, False, default_args)
#agent = Agent(domain = "domain.yml", interpreter = interpreter, action_endpoint = endpoint_conf)
#training_data = agent.load_data("stories.md")
#agent.train(training_data)
#agent.persist("models/dialogue")
