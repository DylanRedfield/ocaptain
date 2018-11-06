from rasa_core.train import online
from rasa_core.interpreter import RasaNLUInterpreter
from rasa_nlu.model import Interpreter
from rasa_core.utils import EndpointConfig
from rasa_core import utils
from rasa_core.agent import Agent

interpreter = RasaNLUInterpreter("models/current/nlu")
endpoint_conf = utils.read_endpoint_config("endpoints.yml", endpoint_type="action_endpoint")
agent = Agent.load("models/dialogue", interpreter = interpreter, action_endpoint = endpoint_conf)
online.serve_agent(agent)
