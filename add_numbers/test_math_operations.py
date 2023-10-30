import unittest
from math_operations import add_numbers

class TestMathOperations(unittest.TestCase):

    def test_add_numbers_positive(self):
        result = add_numbers(3, 5)
        self.assertEqual(result, 8)

    def test_add_numbers_negative(self):
        result = add_numbers(-10, -7)
        self.assertEqual(result, -17)

    def test_add_numbers_mixed(self):
        result = add_numbers(12, -8)
        self.assertEqual(result, 4)

if __name__ == '__main__':
    unittest.main()
