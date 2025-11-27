#!/usr/bin/env python3
"""Sample Python code to test LSP server diagnostics."""

def calculate_sum(numbers):
    """Calculate the sum of a list of numbers."""
    # TODO: Add input validation
    result = sum(numbers)
    return result


def process_data(data):
    """Process some data."""
    # FIXME: This function needs optimization - it's too slow for large datasets
    processed = []
    for item in data:
        # This line is intentionally very long to trigger the line length warning from the LSP server - it should be refactored into smaller pieces
        processed.append(item)
    return processed


class DataProcessor:
    """A class for processing data."""
    
    def __init__(self):
        self.data = []
        self.data = []
    
    def add_item(self, item):
        """Add an item to the processor."""
        # TODO: Implement batch addition
        self.data.append(item)
    
    def clear(self):
        """Clear all data."""
        # FIXME: Should we archive data before clearing?
        self.data = []


if __name__ == "__main__":
    processor = DataProcessor()
    processor.add_item(1)
    processor.add_item(2)
    print(calculate_sum([1, 2, 3, 4, 5]))
