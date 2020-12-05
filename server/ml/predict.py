import pandas as pd
import glob
import csv
from sklearn.externals import joblib

models = glob.glob('saved/*.pkl')

TESTFALG = True

if TESTFALG:
    from pipe import iload_iris_pipe, isave_iris_data
    x_train, y_train, x_test, y_test = iload_iris_pipe()

    for model in models:
        
        clf = joblib.load(model)
        modelname = clf.__class__.__name__
        if hasattr(clf, "predict") and hasattr(clf, 'predict_proba'):
            predicts = clf.predict(x_test)
            predicts_proba = clf.predict_proba(x_test)

            isave_iris_data(predicts, predicts_proba, 'saved/{}.predict'.format(modelname))

def main():
    from pipe import iload_predict_data , isave_predict_data
    data_id, data_features = iload_predict_data()

    for model in models:
        
        clf = joblib.load(model)
        modelname = clf.__class__.__name__

        if hasattr(clf, "predict"):
            _ = clf.predict(data_features)
            save_predict = "saved/{}_predict.csv".format(modelname)
            isave_predict_data(data_id, _, save_predict)

        if hasattr(clf, 'predict_proba'):
            _ = clf.predict_proba(data_features)
            _ = [ 1-i[0] for i in _ ]
            save_predict_proba = "saved/{}_predict_proba.csv".format(modelname)
            isave_predict_data(data_id, _, save_predict_proba)

# if __name__ == '__main__':
    # main()
