## check open
* check_is_open
  - action_check_is_open 

## check time close
* check_time_close
  - action_check_time_close
## Generated Story 5839377152582597870
* inform
    - action_update_order
    - action_ask_next
* inform{"type": "DELIVERY"}
    - utter_ask_address
    - slot{"type": "DELIVERY"}
    - action_update_order
* inform{"address": "2651 deer path , scotch plains"}
    - utter_ask_order_contents
    - slot{"address": "2651 deer path , scotch plains"}
    - action_update_order
* inform{"contents": "two fried rice"}
    - utter_ask_is_all
    - slot{"contents": "two fried rice"}
    - action_update_order
* affirm
    - utter_ask_confirmation_delivery
* affirm
    - utter_after_order
    - action_place_order
    - export

