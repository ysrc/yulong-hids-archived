from sklearn.datasets import load_digits
from sklearn.model_selection import train_test_split

from sklearn.preprocessing import Binarizer,StandardScaler
from sklearn.decomposition import PCA

from sklearn.pipeline import Pipeline, make_pipeline
def iload_digits_pipe():
    digits = load_digits()
    data = digits.data
    lables = digits.target

    pipe = Pipeline([
        ('scale',StandardScaler())
        ('reduce_dim',PCA())
        ])

    data = pipe.fit_transform(data)
    x_train, x_test, y_train, y_test = train_test_split(data, lables, test_size=0.2, random_state=42)

    return x_train, y_train, x_test, y_test


def isave_digits_data(predict, filename, predict_proba=None):
    pass
