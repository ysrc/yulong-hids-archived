from sklearn.pipeline import Pipeline
from sklearn.model_selection import GridSearchCV
from sklearn.feature_selection import SelectKBest,chi2

from sklearn.decomposition import PCA, NMF
from sklearn.svm import SVC



class IGridSVC():
    N_FEATURES_OPTIONS = [2, 4]
    C_OPTIONS = [1, 10, 100, 1000]
    param_grid = [
            {
            'reduce_dim': [PCA(iterated_power=7), NMF()],
            'reduce_dim__n_components': N_FEATURES_OPTIONS,
            'classify__C': C_OPTIONS
        },
            {
            'reduce_dim': [SelectKBest(chi2)],
            'reduce_dim__k': N_FEATURES_OPTIONS,
            'classify__C': C_OPTIONS
        },]

    pipe = Pipeline([
        ('reduce_dim', PCA()),
        ('classify', SVC( kernel="linear", probability=True))
    ])

    def __init__(self):
        self.model = None

    def fit(self,x_train, y_train):
        self.model = GridSearchCV(IGridSVC.pipe, cv=3, n_jobs=-1, param_grid=IGridSVC.param_grid)
        self.model = self.model.fit(x_train,y_train)

    def score(self,x_test,y_test):
        return self.model.score(x_test,y_test)
    
    def predict(self, x_test):
        return self.model.predict(x_test)

    def predict_proba(self, x_test):
        return self.model.predict_proba(x_test)