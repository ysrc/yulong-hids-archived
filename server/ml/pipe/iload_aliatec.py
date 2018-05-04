import pandas as pd
import os

from sklearn.model_selection import train_test_split

train_data_path = 'data/atec_anti_fraud_train.csv'
predict_data_path = 'data/atec_anti_fraud_test_a.csv'

DROPCOLUMS = ["id","label","date"]
# 0 .... 1, 0 is safe / 1 is not safe

def iload_aliatec_pipe():

    if os.path.isfile(train_data_path) and os.path.isfile(predict_data_path):
        print("[√] Path Checked, File Exists")
    else: 
        print("[X] Please Make Sure Your Datasets Was Exists")
        import sys
        sys.exit(1)
        
    data = pd.read_csv(train_data_path)
    data = data.fillna(0)
    unlabeled = data[data['label'] == -1]
    labeled = data[data['label'] != -1]

    train, test = train_test_split(labeled, test_size=0.2, random_state=42)

    cols = [c for c in DROPCOLUMS if c in train.columns]
    x_train = train.drop(cols,axis=1)

    cols = [c for c in DROPCOLUMS if c in test.columns]
    x_test = test.drop(cols,axis=1)

    y_train = train['label']
    y_test = test['label']
    return x_train, y_train, x_test, y_test

def iload_predict_data():
    upload_test = pd.read_csv(predict_data_path)
    upload_test = upload_test.fillna(0)
    upload_id = upload_test['id']
    
    cols = [c for c in DROPCOLUMS if c in upload_test.columns]
    upload_test = upload_test.drop(cols,axis=1)

    return upload_id, upload_test


def isave_predict_data(data_id,predict,filename):
    p = pd.DataFrame(predict,columns=["score"])
    res = pd.concat([data_id,p],axis=1)
    res.to_csv(filename,index=False)
    print("[+] Save Predict Result To {} Sucessful".format(filename))