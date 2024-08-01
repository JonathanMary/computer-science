#define PY_SSIZE_T_CLEAN
#include <Python.h>


static PyObject *cvarint_encode(PyObject *self, PyObject *args)
{
    unsigned long n;
    if (!PyArg_ParseTuple(args, "K", &n))
    {
        return NULL;
    }

    unsigned char out[32];
    int i = 0;

    while (n > 0)
    {
        char part = n & 0x7f;
        n >>= 7;
        part |= (n ? 0x80 : 0x00);
        out[i] = part;
        i += 1;
    }
    return PyBytes_FromStringAndSize((char *)out, i);
}


static PyObject *cvarint_decode(PyObject *self, PyObject *args)
{
    // Method 1, from class
    // const char *varn;
    // if (!PyArg_ParseTuple(args, "y", &varn))
    // {
    //     return NULL;
    // }

    // char b;
    // unsigned long long n = 0;
    // int i, shamt = 0;

    // for (i = 0;; i++)
    // {
    //     b = varn[i];
    //     if (b == 0)
    //     {
    //         break;
    //     }
    //     n |= ((unsigned long long)(b & 0x7f) << shamt);
    //     shamt += 7;
    // }

    // return PyLong_FromUnsignedLongLong(n);

    // Method 2
    PyObject *varn;
    if (!PyArg_ParseTuple(args, "O", &varn))
    {
        return NULL;
    }

    unsigned long n = 0;
    Py_ssize_t varnLength = PyBytes_Size(varn);

    for (Py_ssize_t i = varnLength - 1; i >= 0; i--)
    {
        n <<= 7;
        n |= (PyBytes_AS_STRING(varn)[i] & 0x7f);
    }
    return PyLong_FromUnsignedLong(n);
}


static PyMethodDef CVarintMethods[] = {
    {"encode", cvarint_encode, METH_VARARGS, "Encode an integer as varint."},
    {"decode", cvarint_decode, METH_VARARGS,
     "Decode varint bytes to an integer."},
    {NULL, NULL, 0, NULL}};


static struct PyModuleDef cvarintmodule = {
    PyModuleDef_HEAD_INIT, "cvarint",
    "A C implementation of protobuf varint encoding", -1, CVarintMethods};


PyMODINIT_FUNC PyInit_cvarint(void) { return PyModule_Create(&cvarintmodule); }
