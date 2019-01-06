import json

def generate_output(text, intent, entities, slots, actions):
    correct_slots = list(set(slots) | set(default_slots()))
    output = {'text': text, 'intent': intent, 'entities': entities, 'slots': correct_slots, 'actions': actions}

    return output

def gen_reservation_size_files_all_have():
    output = []
    actions = ['utter_ask_for_number_on_reservation_size']

    output.append(generate_output("I'd like to make a reservation", "make_reservation", [], [], actions))
    return output

def gen_reservation_size_happy():
    output = gen_reservation_size_files_all_have()

    actions = ['action_set_potential_size_slot', 
    'action_brancher_validate_potential_reseservation_size', 
    'action_set_size_slot_from_potential_size_slot',
    'action_clear_potential_size_slot']

    output.append(generate_output("Six", "inform", ["number"], ['size'], actions))
    write_to_file('tests/reservation_size_filling/happy.test', output)

def gen_unrecognized_inform_name():
    output = gen_reservation_size_files_all_have()

    actions = ['utter_unhappy_doing_no_number_recognized']
    output.append(generate_output("Dylan", "inform", ["name"], [], actions))
    write_to_file('tests/reservation_size_filling/gen_unrecognized_inform_name.test', output)

def gen_unrecognized_inform_time():
    output = gen_reservation_size_files_all_have()

    actions = ['utter_unhappy_doing_no_number_recognized']
    output.append(generate_output("Dylan and 7pm", "inform", ["name", "time"], [], actions))
    write_to_file('tests/reservation_size_filling/gen_unrecognized_inform_time.test', output)

def gen_unrecognized_inform_name_and_time():
    output = gen_reservation_size_files_all_have()

    actions = ['utter_unhappy_doing_no_number_recognized']
    output.append(generate_output("Dylan and 7pm", "inform", ["name", "time"], [], actions))
    write_to_file('tests/reservation_size_filling/gen_unrecognized_inform_name_and_time.test', output)

def gen_potential_time_greater_than_max():
    output = gen_reservation_size_files_all_have()

    actions = ['action_set_potential_size_slot', 
            'action_brancher_validate_potential_reservation_size',
            'utter_unhappy_doing_request_customer_call_for_large_parties',
            'action_clear_potential_size_slot']
    output.append(generate_output("twenty one", "inform", ["number"], [], actions))
    write_to_file('tests/reservation_size_filling/gen_potential_time_greater_than_max.test', output)

def gen_potential_time_zero():
    output = gen_reservation_size_files_all_have()

    actions = ['action_set_potential_size_slot', 
            'action_brancher_validate_potential_reservation_size',
            'utter_unhappy_doing_invalid_size_and_ask_for_size_greater_than_zero',
            'action_clear_potential_size_slot']
    output.append(generate_output("0", "inform", ["number"], [], actions))
    write_to_file('tests/reservation_size_filling/gen_potential_time_zero.test', output)

def gen_unrecoginzed_with_name_then_works():
    output = gen_reservation_size_files_all_have()

    actions = ['utter_unhappy_doing_no_number_recognized']
    output.append(generate_output("Dylan", "inform", ["name", "time"], [], actions))

    actions = ['action_set_potential_size_slot', 
    'action_brancher_validate_potential_reseservation_size', 
    'action_set_size_slot_from_potential_size_slot',
    'action_clear_potential_size_slot']

    output.append(generate_output("Six", "inform", ["number"], ['size'], actions))
    write_to_file('tests/reservation_size_filling/happy.test', output)

def gen_two_unrecogized_then_invalid():
    output = gen_reservation_size_files_all_have()

    actions = ['utter_unhappy_doing_no_number_recognized']
    output.append(generate_output("Dylan", "inform", ["name", "time"], [], actions))

    actions = ['action_set_potential_size_slot', 
            'action_brancher_validate_potential_reservation_size',
            'utter_unhappy_doing_invalid_size_and_ask_for_size_greater_than_zero',
            'action_clear_potential_size_slot']
    output.append(generate_output("0", "inform", ["number"], [], actions))
    write_to_file('tests/reservation_size_filling/gen_two_unrecogized_then_invalid.test', output)


def write_to_file(path, output):
    with open(path, 'w') as outfile:
        json.dump(output, outfile)

def default_slots():
    return ['business_id', 'recipient_id', 'recipient_contact']

if __name__ == '__main__':
    gen_two_unrecogized_then_invalid()
    gen_unrecoginzed_with_name_then_works()
    gen_potential_time_zero()
    gen_potential_time_greater_than_max()
    gen_unrecognized_inform_name_and_time()
    gen_unrecognized_inform_time()
    gen_unrecognized_inform_name()
