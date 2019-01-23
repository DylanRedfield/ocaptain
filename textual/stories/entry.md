## entry
* make_reservation
    - action_blank_for_slots
    - slot{"business_id": "MewuHeThW4QJGDxD9tTr"}
    - slot{"recipient_id": "BGeFfREAGGSRqRWrmLNx"}
    - slot{"recipient_contact": "+19084771280"}
    - action_brancher_reservation_slot_filling_base

* make_reservation{"number": 5}
    - action_blank_for_slots
    - slot{"business_id": "MewuHeThW4QJGDxD9tTr"}
    - slot{"recipient_id": "BGeFfREAGGSRqRWrmLNx"}
    - slot{"recipient_contact": "+19084771280"}
    - action_set_potential_size_slot
    - slot{"potential_size": 5}
    - action_brancher_reservation_slot_filling_base

* make_reservation{"time": "xxxx"}
    - action_blank_for_slots
    - slot{"business_id": "MewuHeThW4QJGDxD9tTr"}
    - slot{"recipient_id": "BGeFfREAGGSRqRWrmLNx"}
    - slot{"recipient_contact": "+19084771280"}
    - action_set_temp_times_slot
    - slot{"temp_times": []}
    - action_brancher_reservation_slot_filling_base

* make_reservation{"number": 8, "time": "xxxx"}
    - action_blank_for_slots
    - slot{"business_id": "MewuHeThW4QJGDxD9tTr"}
    - slot{"recipient_id": "BGeFfREAGGSRqRWrmLNx"}
    - slot{"recipient_contact": "+19084771280"}
    - action_set_potential_size_slot
    - slot{"potential_size": 5}
    - action_set_temp_times_slot
    - slot{"temp_times": []}
    - action_brancher_reservation_slot_filling_base

