from rasa_core.slots import Slot
class TypeSlot(Slot):
    type_name = "order_type"
    def feature_dimensionality(self):
        return 2

    def as_feature(self):
        r = [0.0] * self.feature_dimensionality()
        if self.value:
            if self.value == "PICK_UP":
                r[0] = 1.0
            elif self.value == "DELIVERY":
                r[1] = 1.0
        return r
