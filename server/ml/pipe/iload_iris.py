from sklearn.datasets import load_iris
from sklearn.model_selection import train_test_split

import pandas as pd

def iload_iris_pipe():
    iris = load_iris()
    x_train, x_test, y_train, y_test = train_test_split(iris.data,iris.target,test_size=0.2, random_state=42 )
    
    return x_train, y_train, x_test, y_test


def isave_iris_data(predict, predict_proba, filename):
    proba = []
    p1 = pd.DataFrame(predict, columns=["type"])
    print(predict_proba)
    shape = predict_proba.shape
    print(shape)
    proba = [v[k] for k, v in zip(predict, predict_proba)]

    p2 = pd.DataFrame(proba, columns=["proba"])
    res = pd.concat([p1, p2], axis=1)

    res.to_csv(filename, index=False)
    print("[+] Save Predict Result To {} Sucessful".format(filename))
