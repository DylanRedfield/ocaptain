from rasa_core.actions import Action
from rasa_core.events import SlotSet

class CheckTimeClose(Action):
   def name(self):
      # type: () -> Text
      return "action_check_time_close"

   def run(self, dispatcher, tracker, domain):
      # type: (Dispatcher, DialogueStateTracker, Domain) -> List[Event]

      return []

class PlaceOrder(Action):
   def name(self):
      # type: () -> Text
      return "action_place_order"

   def run(self, dispatcher, tracker, domain):
      # type: (Dispatcher, DialogueStateTracker, Domain) -> List[Event]

      return []

class CheckIsOpen(Action):
   def name(self):
      # type: () -> Text
      return "action_check_is_open"

   def run(self, dispatcher, tracker, domain):
      # type: (Dispatcher, DialogueStateTracker, Domain) -> List[Event]

      return []

class StartOrder(Action):
   def name(self):
      # type: () -> Text
      return "action_start_order"

   def run(self, dispatcher, tracker, domain):
      # type: (Dispatcher, DialogueStateTracker, Domain) -> List[Event]

      return []

class CheckIsOpenOnDay(Action):
   def name(self):
      # type: () -> Text
      return "action_check_is_open_on_day"

   def run(self, dispatcher, tracker, domain):
      # type: (Dispatcher, DialogueStateTracker, Domain) -> List[Event]

      return []

class CheckTimeCloseOnDay(Action):
   def name(self):
      # type: () -> Text
      return "action_check_time_close_on_day"

   def run(self, dispatcher, tracker, domain):
      # type: (Dispatcher, DialogueStateTracker, Domain) -> List[Event]

      return []


