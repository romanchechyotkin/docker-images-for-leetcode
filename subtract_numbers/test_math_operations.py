import unittest
from math_operations import subtract_numbers

class TestMathOperations(unittest.TestCase):

    def test_add_numbers_positive(self):
        result = subtract_numbers(3, 5)
        self.assertEqual(result, -2)

    def test_add_numbers_negative(self):
        result = subtract_numbers(-10, -7)
        self.assertEqual(result, -3)

    def test_add_numbers_mixed(self):
        result = subtract_numbers(12, -8)
        self.assertEqual(result, 20)

if __name__ == '__main__':
    unittest.main()
