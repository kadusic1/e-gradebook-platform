def get_difficulty_for_tenant(avg_grade: float) -> str:
    """Convert average grade to category"""
    if avg_grade >= 4.50:
        return "Vrlo lagana"
    elif avg_grade >= 3.50:
        return "Lagana"
    elif avg_grade >= 2.50:
        return "Srednja"
    else:
        return "Teška"


def get_criteria_for_teacher(avg_grade: float) -> str:
    """Convert average grade to category"""
    if avg_grade >= 4.50:
        return "Vrlo lagan"
    elif avg_grade >= 3.50:
        return "Lagan"
    elif avg_grade >= 2.50:
        return "Srednji"
    else:
        return "Težak"


def get_success_category(avg_grade: float) -> str:
    """Convert average grade to success category"""
    if avg_grade >= 4.50:
        return "Odličan"
    elif avg_grade >= 4.00:
        return "Vrlo dobar"
    elif avg_grade >= 3.00:
        return "Dobar"
    elif avg_grade >= 2.00:
        return "Dovoljan"
    else:
        return "Nedovoljan"


def safe_float_str(x):
    # Force plain float, never scientific notation
    # 12 decimals is usually enough for embeddings
    return format(float(x), ".12f").rstrip("0").rstrip(".")
