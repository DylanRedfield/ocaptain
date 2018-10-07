from rasa_core.agent import Agent
from rasa_core.interpreter import RasaNLUInterpreter
from rasa_core import utils
from rasa_core.train import train_dialogue_model
#from rasa_core.utils import AvailableEndpoints
from rasa_core.run import AvailableEndpoints

interpreter = RasaNLUInterpreter("models/current/nlu")

endpoint_conf = utils.read_endpoint_config("endpoints.yml", endpoint_type="action_endpoint")
endpoints = AvailableEndpoints(nlg = None, nlu = None, action = endpoint_conf, model = None)
agent = train_dialogue_model("domain.yml", "stories.md", "models/dialogue", endpoints)
#agent = Agent(domain = "domain.yml", interpreter = interpreter, action_endpoint = endpoint_conf)
#training_data = agent.load_data("stories.md")
#agent.train(training_data)
#agent.persist("models/dialogue")
