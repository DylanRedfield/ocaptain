## entry
* make_reservation
    - action_blank_for_slots
    - slot{"business_id": "MewuHeThW4QJGDxD9tTr"}
    - slot{"recipient_id": "BGeFfREAGGSRqRWrmLNx"}
    - slot{"recipient_contact": "+19084771280"}
    - action_brancher_reservation_slot_filling_base

## entry2
* make_reservation{"number": 5}
    - action_blank_for_slots
    - slot{"business_id": "MewuHeThW4QJGDxD9tTr"}
    - slot{"recipient_id": "BGeFfREAGGSRqRWrmLNx"}
    - slot{"recipient_contact": "+19084771280"}
    - action_set_potential_size_slot
    - slot{"potential_size": 5}
    - action_brancher_reservation_slot_filling_base

## entry3
* make_reservation{"time": "xxxx"}
    - action_blank_for_slots
    - slot{"business_id": "MewuHeThW4QJGDxD9tTr"}
    - slot{"recipient_id": "BGeFfREAGGSRqRWrmLNx"}
    - slot{"recipient_contact": "+19084771280"}
    - action_set_temp_times_slot
    - slot{"temp_times": []}
    - action_brancher_reservation_slot_filling_base

## entry4
* make_reservation{"number": 8, "time": "xxxx"}
    - action_blank_for_slots
    - slot{"business_id": "MewuHeThW4QJGDxD9tTr"}
    - slot{"recipient_id": "BGeFfREAGGSRqRWrmLNx"}
    - slot{"recipient_contact": "+19084771280"}
    - action_set_potential_size_slot
    - slot{"potential_size": 5}
    - action_set_temp_times_slot
    - slot{"temp_times": ["xxx"]}
    - action_brancher_reservation_slot_filling_base

## question_time
* question_time
    - action_utter_answer_time

## question_serve_alcohol
* question_serve_alcohol
    - utter_answer_serve_alcohol

## question_dress_code
* question_dress_code
    - utter_answer_dress_code

## question_gluten_free
* question_gluten_free
    - utter_answer_gluten_free

## question_vegetarian
* question_vegetarian
    - utter_answer_vegetarian

## question_vegan
* question_vegan
    - utter_answer_vegan

## question_kosher
* question_kosher
    - utter_answer_kosher

## question_large_group
* question_large_group
    - utter_answer_large_group

## question_kids_menu
* question_kids_menu
    - utter_answer_kids_menu

## question_byob
* question_bring_your_own_alcohol
    - utter_answer_bring_your_own_alcohol

## question_parking
* question_parking
    - utter_answer_parking

## question_handicap_accessible
* question_handicap_accessible
    - utter_answer_handicap_accessible

## question_corking_fee
* question_corking_fee
    - utter_answer_corking_fee

## question_polar_reservation
* question_polar_reservation
    - utter_answer_polar_reservation

