import json

def generate_output(intent, entities, slots, actions):
    correct_slots = list(set(slots) | set(default_slots()))
    output = {'intent': intent, 'entities': entities, 'slots': correct_slots, 'actions': actions}

    return output

def gen_1():
    output = []
    # 'Can I get a table for two?'
    output.append(generate_output("inform", ["number"], ['size'], ['action_set_size',
        'utter_ask_time','action_listen']))
    # 'tomorrow at 7pm'
    actions = ['action_set_scheduled_time', 'action_check_reservation_datetime','action_ask_if_any_similar_times_work',
    'action_listen'] 
    output.append(generate_output("inform", ["time"], ['size', 'scheduled_time'], actions))
    with open('tests/output/1.output', 'w') as outfile:
        json.dump(output, outfile)

def default_slots():
    return ['business_id', 'recipient_id', 'recipient_contact']

if __name__ == '__main__':
    gen_1()
