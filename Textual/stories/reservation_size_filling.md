## checkpoint_to_ask_for_potential_size
> checkpoint_to_ask_for_potential_size
    - utter_ask_for_number_on_reservation_size
> checkpoint_to_fill_potential_size

## checkpoint_to_fill_potential_size_inform_number
> checkpoint_to_fill_potential_size
* inform{"number": 2}
    - action_set_potential_size_slot
    - slot{"potential_size": 1}
> checkpoint_with_potential_size_to_validate_and_fill_size

## checkpoint_to_fill_potential_size_inform_non_number
> checkpoint_to_fill_potential_size
* inform{"name": "dylan"}
    - utter_unhappy_doing_no_number_recognized
> checkpoint_to_fill_potential_size

## checkpoint_to_fill_potential_size_inform_none
> checkpoint_to_fill_potential_size
* inform
    - utter_unhappy_doing_no_number_recognized
> checkpoint_to_fill_potential_size

## checkpoint_with_potential_size_to_validate_and_fill_size_where_potential_size_is_valid
> checkpoint_with_potential_size_to_validate_and_fill_size
    - action_brancher_validate_reservation_potential_size
    - slot{"size": 5}
    - action_clear_potential_size_slot

## checkpoint_with_potential_size_to_validate_and_fill_size_where_potential_size_is_zero
> checkpoint_with_potential_size_to_validate_and_fill_size
    - action_brancher_validate_reservation_potential_size
    - followup{"name": "utter_unhappy_doing_invalid_size_AND_ask_for_size_greater_than_zero"}
    - utter_unhappy_doing_invalid_size_AND_ask_for_size_greater_than_zero
    - action_clear_potential_size_slot
> checkpoint_to_fill_potential_size

## checkpoint_with_potential_size_to_validate_and_fill_size_where_potential_size_is_greater_than_max
> checkpoint_with_potential_size_to_validate_and_fill_size
    - action_brancher_validate_reservation_potential_size
    - followup{"name": "utter_unhappy_doing_request_customer_call_for_large_parties"}
    - utter_unhappy_doing_request_customer_call_for_large_parties
    - action_clear_potential_size_slot
    - action_need_employee

## checkpoint_with_potential_size_to_validate_and_fill_size_where_potential_size_is_null
> checkpoint_with_potential_size_to_validate_and_fill_size
    - action_brancher_validate_reservation_potential_size
    - followup{"name": "utter_ask_for_number_on_reservation_size"}
    - utter_ask_for_number_on_reservation_size
    - action_clear_potential_size_slot
> checkpoint_to_fill_potential_size

