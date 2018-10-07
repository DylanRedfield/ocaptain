import logging
import time

from flask import Blueprint, request, jsonify

from twilio.rest import Client

from rasa_core.channels import UserMessage, OutputChannel
from rasa_core.channels import InputChannel
import firebase_admin
from firebase_admin import credentials
from firebase_admin import firestore
import requests
import threading

class TextualOutput(OutputChannel):

    @classmethod
    def name(cls):
        return "textual"
    
    def __init__(self, req, db, sms_client):
        self.req = req
        self.sms_client = sms_client
        self.db = db


    def send_text_message(self, recipient_id, text):
        contact = self.req["recipient"]["Contact"]
        business_phone = self.req["business"]["PhoneNumber"]
        business_id = self.req["business"]["Id"]
        recipient_id = self.req["recipient"]["Id"]
        self.sms_client.messages.create(body=text, from_= business_phone, to = contact)

        message = {u'content': text, u'didBotCreate':True, u'hasBusinessRead':False,u'isBusinessSender':True,
                u'recipientId': recipient_id, u'timeSent':round(time.time() * 1000)} 


        self.db.collection(u'businesses').document(business_id).collection(u'messages').add(message)
        recipient_ref = self.db.collection(u'businesses').document(business_id).collection(u'recipients').document(recipient_id)

        recipient_ref.update({'recentMessage': message})


class TextualInput(InputChannel):

    def __init__(self, agent):
        self.agent = agent

    @classmethod
    def name(cls):
        return "textual"

    def blueprint(self, on_new_message):
        cred = credentials.Certificate('firebase-config.json')
        firebase_admin.initialize_app(cred)
        db = firestore.client()

        account_sid = 'AC9dfbda388f3ee10353bbc001694f5c27'
        auth_token = 'e3429e06cc27740f1c859d2bfc9964ae'
        sms_client = Client(account_sid, auth_token)

        textual_webhook = Blueprint('textual_webhook', __name__)

        @textual_webhook.route("/", methods=['GET'])
        def health():
            return jsonify({"status":"ok"})

        @textual_webhook.route("/webhook", methods=['POST'])
        def message():
            req = request.get_json()

            text = req["message"]["Content"]
            sender = req["recipient"]["Id"]

            tracker = self.agent.tracker_store.get_or_create_tracker(sender)

            if tracker.slots['business_id'] == "" or tracker.slots['recipient_contact'] == "":
                tracker._set_slot('business_id', req["business"]["Id"])
                tracker._set_slot('recipient_contact', req["recipient"]["Contact"])
                print("In the method")
                print(req["recipient"]["Contact"])
                self.agent.tracker_store.save(tracker)

            out_channel = TextualOutput(req, db, sms_client)

            user = UserMessage(text, output_channel = out_channel, sender_id = sender)

            print(user)
            on_new_message(user)
            return "success"
        return textual_webhook

