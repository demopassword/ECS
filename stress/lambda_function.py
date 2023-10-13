import json
import boto3
import time

cf = boto3.client('cloudfront')
code_pipeline = boto3.client('codepipeline')

DISTRIBUTION_ID = "EX46VUFBQTI6A"

def create_invalidation():
    res = cf.create_invalidation(
        DistributionId=DISTRIBUTION_ID,
        InvalidationBatch={
            'Paths':{
                'Quantity': 1,
                'Items': [
                    '/v1/stress*',
                ]
            },
            'CallerReference': str(time.time())
        }
    )
    invalidation_id = res['Invalidation']['Id']
    return invalidation_id

def lambda_handler(event, context):
    id = create_invalidation()
    job = event['CodePipeline.job']['id']
    
    if id:
        return code_pipeline.put_job_success_result(jobId=job)
    else:
        return code_pipeline.put_job_failure_result(jobId=job, failureDetails={'message': message, 'type': 'JobFailed'})
