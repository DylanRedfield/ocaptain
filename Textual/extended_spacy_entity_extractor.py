from __future__ import absolute_import
from __future__ import division
from __future__ import print_function
from __future__ import unicode_literals

import typing
from typing import Any
from typing import Dict
from typing import List
from typing import Text

from rasa_nlu.extractors import EntityExtractor
from rasa_nlu.training_data import Message

if typing.TYPE_CHECKING:
    from spacy.tokens.doc import Doc


class ExtendedSpacyEntityExtractor(EntityExtractor):
    name = "ner_extended_spacy"

    provides = ["entities"]

    requires = ["spacy_nlp"]

    def __init__(self, component_config=None):
        super(ExtendedSpacyEntityExtractor, self).__init__(component_config)
        self.ordinals = []

    def process(self, message, **kwargs):
        # type: (Message, **Any) -> None

        # can't use the existing doc here (spacy_doc on the message)
        # because tokens are lower cased which is bad for NER
        spacy_nlp = kwargs.get("spacy_nlp", None)
        doc = spacy_nlp(message.text)
        entities = self.extract_entities(doc)

        self.ordinals = []
        for entity in entities:
            if "ORDINAL" == entity['entity']:
                self.ordinals.append(entity['value'])

        entities = list(filter(lambda x: "ORDINAL" != entity['entity'], entities))


        extracted = self.add_extractor_name(entities)
        message.set("entities",
                    message.get("entities", []) + extracted,
                    add_to_output=True)
        message.set("ordinals", self.ordinals)



    @staticmethod
    def extract_entities(doc):
        # type: (Doc) -> List[Dict[Text, Any]]

        entities = [
            {
                "entity": ent.label_,
                "value": ent.text,
                "start": ent.start_char,
                "confidence": None,
                "end": ent.end_char
            }
            for ent in doc.ents]
        return entities
