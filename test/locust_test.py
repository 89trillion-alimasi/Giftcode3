from locust import HttpUser, between, task

import locust_test


class QuickstartUser(HttpUser):
    wait_time = between(3, 5)


    #创建接口的测试
    @task(3)
    def createGiftCode(self):
        data={
            "type": 3,
            "receiving_user": "alms",
            "valid_period": "2021-06-29 17:48:00",
            "create_user": "admin",
            "available_times": 2,
            "description": "测试type1",
            "gift_packages": [
                {
                    "name": "金币",
                    "num": 10
                },
                {
                    "name": "钻石",
                    "num": 20
                }
            ]
        }


        self.client.post("/create_gift_code",json=data)


    #查询接口的测试
    @task(7)
    def queryGiftCode(self):


        self.client.get("/query_gift_code?code=HVALJMX9")


    #验证码测试
    @task(2)
    def verifyGiftCode(self):
        data={
                "user": "alms",
                "code": "HVALJMX9"
            }

        self.client.post("/verify_gift_code",json=data)