## Generated Story -2653337070988492228
* make_reservation{"number": 2}
    - action_blank_for_slots
    - slot{"business_id": "MewuHeThW4QJGDxD9tTr"}
    - slot{"recipient_id": "BGeFfREAGGSRqRWrmLNx"}
    - slot{"recipient_contact": "+19084771280"}
    - action_set_potential_size_slot
    - slot{"potential_size": 2}
    - action_brancher_reservation_slot_filling_base
    - followup{"name": "action_checkpoint_with_potential_size_to_validate_and_fill_size"}
    - action_checkpoint_with_potential_size_to_validate_and_fill_size
    - action_brancher_validate_reservation_potential_size
    - slot{"size": 2}
    - followup{"name": "action_blank_alert_size_slot_set"}
    - action_blank_alert_size_slot_set
    - action_clear_potential_size_slot
    - slot{"potential_size": null}
    - action_brancher_reservation_slot_filling_base
    - followup{"name": "utter_ask_for_time_for_potential_reservation"}
    - utter_ask_for_time_for_potential_reservation

