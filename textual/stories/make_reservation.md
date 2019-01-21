## make reservatio no input
* make_reservation
    - action_blank_for_slots
    - slot{"business_id": "MewuHeThW4QJGDxD9tTr"}
    - slot{"recipient_id": "BGeFfREAGGSRqRWrmLNx"}
    - slot{"recipient_contact": "+19084771280"}
> checkpoint_to_ask_for_potential_size

## Generated Story -6719338366551963822
* make_reservation{"number": 2}
    - action_blank_for_slots
    - slot{"business_id": "MewuHeThW4QJGDxD9tTr"}
    - slot{"recipient_id": "BGeFfREAGGSRqRWrmLNx"}
    - slot{"recipient_contact": "+19084771280"}
    - action_set_potential_size_slot
    - slot{"potential_size": 2}
> checkpoint_with_potential_size_to_validate_and_fill_size

