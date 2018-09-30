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

class TextualOutput(OutputChannel):

    @classmethod
    def name(cls):
        return "textual"
    
    def __init__(self, req, firebase_client, sms_client):
        self.req = req
        self.firebase_client = firebase_client
        self.sms_client = sms_client

    def send_text_message(self, recipient_id, text):
        sms_client.messages.create(body=text, from_=req["recipient"]["contact"],to=req["business"]["phoneNumber"])
        ref = firebase_client.collection(u'bussinesses').docuemnt(req["business"]["id"]).collection("messages")

        message = {'content':text, 'didBotCreate':True, 'hasBusinessRead':False,'isBusinessSender':True, 'recipientId':req["recipient"]["id"], 'timeSent':time.time()} 
        ref.add(message)

        recipient_ref = firebase_client.collection(u'bussinesses').docuemnt(req["business"]["id"]).collection("recipients").document(req["recipient"]["id"])

        recipient_ref.update({'recentMessage': message})

class TextualInput(InputChannel):

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
            req = request.form.to_dict()
            requests.post("http://localhost:5005/conversations/{}/tracker/events".format(req["recipient"]["id"]),
                    data = {'event':'slot','name':'business_id', 'value': req["business"]["id"]
            requests.post("http://localhost:5005/conversations/{}/tracker/events".format(req["recipient"]["id"]),
                    data = {'event':'slot','name':'recipient_contact', 'value': req["recipient"]["contact"]

            out_channel = TextualOutput(request.form.to_dict(), db, sms_client)

            on_new_message(UserMessage(request.form.to_dict()["message"]["content"], out_channel, sender))
            return "success"
        return textual_webhook
