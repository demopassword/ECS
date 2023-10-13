import json
import boto3
import time

cf = boto3.client('cloudfront')

DISTRIBUTION_ID = "EX46VUFBQTI6A"

def lambda_handler(event, context):
    cf.create_invalidation(
    DistributionId=DISTRIBUTION_ID,
    InvalidationBatch={
        'Paths':{
            'Quantity': 2,
            'Items': [
                '/',
                '/static/index.html'
            ]
        },
        'CallerReference': str(time.time())
    })
    return "Success"
