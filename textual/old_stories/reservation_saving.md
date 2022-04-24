## action_brancher_to_save_reservation_0
    - action_brancher_to_save_new_reservation
    - followup{"name": "action_brancher_reservation_slot_filling_base"}
    - action_brancher_reservation_slot_filling_base

## action_brancher_to_save_reservation_1
    - action_brancher_to_save_new_reservation
    - followup{"name": "action_need_employee_because_error_saving"}
    - action_need_employee_because_error_saving

## action_brancher_to_save_reservation_2
    - action_brancher_to_save_new_reservation
    - followup{"name": "action_utter_post_reservation_save_AND_ask_for_next_general_request"}
    - action_utter_post_reservation_save_AND_ask_for_next_general_request
    - action_clear_name_slot
    - slot{"name": null}
    - action_clear_scheduled_time_slot
    - slot{"scheduled_time": null}
    - action_clear_size_slot
    - slot{"size": null}
> checkpoint_general_return
