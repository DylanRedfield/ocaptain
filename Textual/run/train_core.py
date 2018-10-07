from rasa_core.agent import Agent
from rasa_core.interpreter import RasaNLUInterpreter
from rasa_core import utils

interpreter = RasaNLUInterpreter("../models/current/nlu")

endpoint_conf = utils.read_endpoint_config("../endpoints.yml", endpoint_type="action_endpoint")

agent = Agent(domain = "../domain.yml", interpreter = interpreter, action_endpoint = endpoint_conf)
training_data = agent.load_data("../rasa_dataset_training.json")
agent.train(training_data)
agent.persist("models/dialogue")
