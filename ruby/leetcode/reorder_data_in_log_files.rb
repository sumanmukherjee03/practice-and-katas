# @param {String[]} logs
# @return {String[]}
def reorder_log_files(logs)
    letter_logs = []
    digi_logs = []
    logs.each do |l|
        x = l.split(' ')
        if x[1].match(/\d+\s*/)
            digi_logs << l

        else
            letter_logs << l
        end
    end
    return letter_logs.sort do |a, b|
        x = a.split(' ')[1..-1].join(' ').chars
        y = b.split(' ')[1..-1].join(' ').chars
        res = 0
        if x == y
            res = a.split(' ')[0].chars <=> b.split(' ')[0].chars
        else
            res = x <=> y
        end
        res
    end.concat(digi_logs)
end
