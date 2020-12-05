# Code Snippetss


* 动态创建类的函数

```python
def set_class_func(CLASS, FUNC):
    try:
        funclist = [i for i in dir(FUNC) if callable(getattr(FUNC,i)) and not i.startswith('_') ]
        for func in funclist:
            setattr(CLASS,str(func),classmethod(func))    
    except Exception as e:
        import sys
        print("Error When Create Class Method From {} TO {}".format(CLASS,FUNC))
        sys.exit(1)
```

* 文件夹路径
```python
import os
print(os.path.dirname(os.path.abspath("__file__")))
```

* print without new line
```python
print('Loading Data...',end='')
print('Done)
```

* AttributeError: 'LinearSVC' object has no attribute 'predict_proba', https://stackoverflow.com/questions/35076586/linearsvc-vs-svckernel-linear-conflicting-arguments
```python

LinearSVC(probability=True) # 提示没有该属性
#Similar to SVC with parameter kernel=’linear’, but implemented in terms of liblinear rather than libsvm, so it has more flexibility in  the choice of penalties and loss functions and should scale better to large numbers of samples.
# http://scikit-learn.org/stable/modules/svm.html
SVC( kernel="linear", probability=True)

```