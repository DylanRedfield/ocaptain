from rasa_core.agent import Agent

from rasa_core.interpreter import RasaNLUInterpreter
from textual_channel import TextualInput
from rasa_core.channels.console import CmdlineInput
from rasa_core.channels.channel import RestInput
from rasa_core import utils
from pathlib import Path
import unittest
import json

class TestCore(unittest.TestCase):
    def setUp(self):
        print("-- Start interpreter --")
        interpreter = RasaNLUInterpreter("models/current/nlu")

        print("-- Start Load --")
        endpoint_conf = utils.read_endpoint_config("endpoints.yml", endpoint_type="action_endpoint")
        self.agent = Agent.load("models/dialogue", interpreter = interpreter, action_endpoint = endpoint_conf)
        self.sender_id = "+19084771280"

    def test_all_input_files(self):
        pathlist = Path("tests/input").glob('*.input')
        
        for path in pathlist:
            input_file = path.open()
            output_file = open(str(path)[:-5] + "output")

            print(input_file.read())

            input_file.close()
            output_file.close()

    def helper_test_file(self, input_file, output_file):
        outputs = json.load(output_file)

        for i, line in enumerate(input_file):
            self.agent.handle_text(line, sender_id = self.sender_id)
            prediction = self.agent.predict_next(self.sender_id)

            output = outputs[i]

            # Correct intent
            # TODO allow for several
            self.assertEqual(prediction['tracker']["latest_message"]['intent']['name'], output['intent'])

            # Correct entities
            goal_entities = output['entities']
            result_entities = prediction['tracker']['latest_message']['entities']

            contains_all = True
            for entity in goal_entities:
                contains = False
                for result_ent in result_entities:
                    if result_ent['entity'] == entity:
                        contains = True
                if not contains:
                    contains_all = False
            self.assertTrue(contains_all)

            # TODO correct slots
            # TODO correct actions




if __name__ == '__main__':
    unittest.main()
