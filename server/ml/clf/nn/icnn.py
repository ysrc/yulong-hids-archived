import os

from sklearn.model_selection import GridSearchCV

from keras.wrappers.scikit_learn import KerasClassifier
from keras.utils import to_categorical
# Keras Model/Layers

from keras.models import Model,Sequential
from keras.layers import Input, Dense, LSTM, Activation, Conv2D, \
                         MaxPool2D, Dropout, Flatten, Embedding,Reshape,Concatenate,\
                         TimeDistributed, AveragePooling1D
from keras.callbacks import ModelCheckpoint
from keras.optimizers import Adam, Adadelta

from keras.losses import categorical_crossentropy

class ICNN():
    
    # filter_sizes = [3, 4, 5]
    # num_filters = 128
    # epochs = 200
    # batch_size = 64
    num_classes = 4
    
    param_grid = {
        'clf__optimizer': ['rmsprop', 'adam', 'adagrad'],
        'clf__epochs': [200, 300, 400, 700, 1000],
        'clf__batch_size': [32, 64, 128],
        'clf__dropout': [0.1, 0.2, 0.3, 0.4, 0.5],
        'clf__kernel_initializer': ['he_normal', 'glorot_uniform', 'normal', 'uniform']
    }
    
    # pipline = Pipeline([
    #     # ('preprocess_step1',None),
    #     # ('preprocess_step2',None),
    #     # ('preprocess_step3',None)
    #     # ('clf', keras_clf)
    # ])

    def __init__(self,kernel_initializer='he_normal',optimizer='adam',activation='relu',loss='binary_crossentropy',dropout=0.5):
        
        self.kernel_initializer = kernel_initializer
        self.optimizer = optimizer
        self.activation = activation
        self.dropout = dropout
        self.loss = loss
        self.model = None

    
    def creat_model(self):
        model = ""
        return model

    def search_model(self):
        pass
    
    def fit(self,x_train,y_train):
        # self.model.fit(x_train, y_train, batch_size=batch_size, epochs=epochs, verbose=1, callbacks=[checkpoint], validation_data=(x_test, y_test))
        # self.model.save("{}/model.h5".format(self.name))
        # checkpoint = ModelCheckpoint('{}/weights.{epoch:03d}-{val_acc:.4f}.hdf5'.format(self.name), monitor='val_acc', verbose=1, save_best_only=True, mode='auto')

        self.model = KerasClassifier(build_fn=self.creat_model)

        # self.model = GridSearchCV(KerasClassifier(build_fn=self.creat_model),\
        #                 cv=3, param_grid=ICNN.param_grid)

        self.model.fit(x_train,y_train)

    def score(self, x_test,y_test):
        # self.model.score(x_test,y_test)
        return 1