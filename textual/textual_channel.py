import logging
import time

from flask import Blueprint, request, jsonify

from twilio.rest import Client

from rasa.core.channels.channel import UserMessage, OutputChannel, InputChannel
from rasa.shared.core.events import SlotSet
import firebase_admin
from firebase_admin import credentials
from firebase_admin import firestore
import requests
import threading
from sanic import Blueprint, response
from sanic.request import Request
from sanic.response import HTTPResponse
from typing import Text, Callable, Awaitable

class TextualOutput(OutputChannel):

    @classmethod
    def name(cls):
        return "textual"
    
    def __init__(self, req, db, sms_client):
        self.req = req
        self.sms_client = sms_client
        self.db = db


    def send_text_message(self, recipient_id, text):
        print("--- Output Text ---")
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


    @classmethod
    def name(cls) -> Text:
        return "textual"

    def blueprint(self, on_new_message: Callable[[UserMessage], Awaitable[None]]) -> Blueprint:
        cred = credentials.Certificate('firebase-config.json')
        firebase_admin.initialize_app(cred)
        db = firestore.client()

        account_sid = 'AC9dfbda388f3ee10353bbc001694f5c27'
        auth_token = 'e3429e06cc27740f1c859d2bfc9964ae'
        sms_client = Client(account_sid, auth_token)

        textual_webhook = Blueprint('textual_webhook', __name__)

        @textual_webhook.route("/", methods=['GET'])
        def health(request: Request) -> HTTPResponse:
            return jsonify({"status":"ok"})

        @textual_webhook.route("/webhook", methods=['POST'])
        async def receive(request: Request) -> HTTPResponse:
            print("--- Input Text ---")
            req = request.json

            text = req["message"]["Content"]
            sender = req["recipient"]["Id"]

            out_channel = TextualOutput(req, db, sms_client)

            user = UserMessage(text, output_channel = out_channel, sender_id = sender)
            await on_new_message(user)

            return response.json("success")
            print("blueprint test")
        return textual_webhook

